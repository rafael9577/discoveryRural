import React, { useState, useEffect } from "react";
import { View, Text, TouchableOpacity, Image, ScrollView, Modal, ActivityIndicator } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { Ionicons } from "@expo/vector-icons";

import styled from "./styled";
import { detalheAnimal } from "../../../../services/API";

const Linha = ({ rotulo, valor }) =>
  valor ? (
    <View style={styled.linhaInformacao2}>
      <Text style={styled.textoInformacao2}>{rotulo}</Text>
      <Text style={styled.valorInformacao2}>{valor}</Text>
    </View>
  ) : null;

export default ({ navigation, route }) => {
  const { id, resumo } = route.params || {};

  const [animal, setAnimal] = useState(resumo || null); // já mostra o resumo enquanto carrega
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [modalVisivel, setModalVisivel] = useState(false);

  const carregar = async () => {
    try {
      setLoading(true);
      setError(null);
      const dados = await detalheAnimal(id);
      setAnimal(dados);
    } catch (err) {
      console.error("Erro ao carregar detalhe:", err);
      setError("Não foi possível carregar os detalhes.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    carregar();
  }, [id]);

    const a = animal || {};
  const nome = a.nome_popular || "Sem nome";
  const cientifico = a.nome_cientifico;
  const peso = a.peso_medio;       // a API ainda não envia
  const tamanho = a.tamanho_medio; // a API ainda não envia
  const descricao = a.breve_descricao;
  const habito = a.habito;
  const conservacao = a.estado_conservacao;
  const raridade = a.raridade != null ? `${a.raridade}/5` : null;
  const peconhento = a.peconhento == null ? null : a.peconhento ? "Sim" : "Não";
  const vocalizacao = a.vocalizacao_postura;
  const aviso = a.aviso_seguranca;
  const protocolo = a.protocolo_incidente;
  const resgate = a.contatos_resgate;
  const urlImagem = a.imagem_url || "https://via.placeholder.com/300x180?text=Sem+Imagem";

  
  return (
    <SafeAreaView style={styled.safeArea}>
      <View style={styled.containerCabecalho}>
        <TouchableOpacity onPress={() => navigation.goBack()}>
          <Ionicons name="arrow-back" size={28} color="#FFFFFF" />
        </TouchableOpacity>
        <Text style={styled.tituloCabecalho}>Informações do animal</Text>
      </View>

      {error && !animal ? (
        <View style={{ flex: 1, justifyContent: "center", alignItems: "center" }}>
          <Text style={{ color: "#FFFFFF", marginBottom: 12 }}>{error}</Text>
          <TouchableOpacity onPress={carregar}>
            <Text style={{ color: "#A7F3D0", fontWeight: "bold" }}>Tentar novamente</Text>
          </TouchableOpacity>
        </View>
      ) : !animal ? (
        <View style={{ flex: 1, justifyContent: "center", alignItems: "center" }}>
          <ActivityIndicator size="large" color="#FFFFFF" />
        </View>
      ) : (
        <ScrollView style={styled.suporteCartao} showsVerticalScrollIndicator={false}>
          <View style={styled.containerPrincipal}>
            <View style={styled.blocoSuperior}>
              <TouchableOpacity onPress={() => setModalVisivel(true)} activeOpacity={0.8}>
                <Image source={{ uri: urlImagem }} style={styled.imagemAnimal} />
              </TouchableOpacity>

              <View style={styled.divisorVertical} />

              <View style={styled.containerInformacoes}>
                <Text style={styled.nomeAnimal}>{nome}</Text>
                <View style={styled.divisorHorizontal} />
                {cientifico ? <Text style={styled.nomeCientifico}>({cientifico})</Text> : null}

                {peso ? (
                  <View style={styled.linhaInformacao}>
                    <Text style={styled.textoInformacao}>Peso médio:</Text>
                    <Text style={styled.valorInformacao}>{peso}</Text>
                  </View>
                ) : null}

                {tamanho ? (
                  <View style={styled.linhaInformacao}>
                    <Text style={styled.textoInformacao}>Tamanho médio:</Text>
                    <Text style={styled.valorInformacao}>{tamanho}</Text>
                  </View>
                ) : null}
              </View>
            </View>

            <View style={styled.blocoDescricao}>
              {descricao ? <Text style={styled.tituloDescricao2}>{descricao}</Text> : null}

              <Linha rotulo="Hábito" valor={habito} />
              <Linha rotulo="Conservação" valor={conservacao} />
              <Linha rotulo="Raridade" valor={raridade} />

              {vocalizacao ? (
                <>
                  <Text style={styled.tituloDescricao}>Vocalização / Postura: </Text>
                  <Text style={styled.valorInformacao2}>{vocalizacao}</Text>
                </>
              ) : null}

              {aviso || protocolo || resgate ? (
                <View style={styled.containerAviso}>
                  <Ionicons name="warning" size={20} color="#D84315" style={styled.iconeAviso} />
                  {aviso ? (
                    <>
                      <Text style={styled.tituloAviso}>Aviso de Segurança: </Text>
                      <Text style={styled.textoAviso}>{aviso}</Text>
                    </>
                  ) : null}
                  {protocolo ? (
                    <>
                      <Text style={styled.tituloAviso}>Protocolo: </Text>
                      <Text style={styled.textoAviso}>{protocolo}</Text>
                    </>
                  ) : null}
                  {resgate ? (
                    <View style={styled.cabecalhoAviso}>
                      <Text style={styled.tituloAviso}>Resgate: </Text>
                      <Text style={styled.textoAviso}>{resgate}</Text>
                    </View>
                  ) : null}
                </View>
              ) : null}
            </View>
          </View>
        </ScrollView>
      )}

      <Image
        style={{ width: 96, height: 96, alignSelf: "center", marginTop: 15, marginBottom: 15 }}
        source={{ uri: "https://github.com/rafael9577/discoveryRural/blob/main/front-end/assets/logoDiscoveryRuralWhite.png?raw=true" }}
      />

      <Modal
        visible={modalVisivel}
        transparent={true}
        animationType="fade"
        onRequestClose={() => setModalVisivel(false)}
      >
        <TouchableOpacity
          style={styled.modalContainer}
          activeOpacity={1}
          onPress={() => setModalVisivel(false)}
        >
          <TouchableOpacity style={styled.botaoFecharModal} onPress={() => setModalVisivel(false)}>
            <Ionicons name="close" size={32} color="#FFFFFF" />
          </TouchableOpacity>
          <Image source={{ uri: urlImagem }} style={styled.imagemModal} resizeMode="contain" />
        </TouchableOpacity>
      </Modal>
    </SafeAreaView>
  );
};